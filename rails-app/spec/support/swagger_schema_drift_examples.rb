RSpec.shared_examples 'swagger_schema_drift_detection' do |model_class, options = {}|
  let(:schema) { model_class.to_swagger_schema }
  let(:properties) { schema[:properties].keys }
  let(:ignored_attributes) { options[:ignored_attributes] || [] }

  it 'includes all attributes in the swagger schema' do
    # This test ensures that fields documented in Swagger are actually handled by the model
    properties.each do |property|
      # Map camelCase to snake_case for the attribute check
      attr_name = property.to_s.underscore.to_sym

      expect(model_class.new).to respond_to(attr_name), "Expected #{model_class} to have attribute :#{attr_name} (from swagger property :#{property})"
    end
  end

  it 'maps all swagger properties in initialize' do
    # Create a hash with all properties from the schema
    params = properties.each_with_object({}) do |prop, hash|
      if prop == :enrollmentDetails
         hash[prop] = [ { schoolName: "Test School" } ]
      elsif prop == :metadata
         hash[prop] = { durationMs: 100 }
      elsif prop == :address
         hash[prop] = { street1: "123 Main St", city: "Arlington", state: "VA", postalCode: "22202", country: "USA" }
      else
         hash[prop] = "test_value_for_#{prop}"
      end
    end

    instance = model_class.new(params)

    properties.each do |property|
      attr_name = property.to_s.underscore
      value = instance.send(attr_name)
      expected_value = params[property]
      expect(value).to eq(expected_value), "Expected attribute :#{attr_name} to be set from swagger property :#{property}"
    end
  end

  it 'warns if there are attributes not documented in swagger' do
    # This is the "drift" check in the other direction - code has it but swagger doesn't
    # Get all attr_accessors (public methods that have a setter equivalent)
    model_instance = model_class.new
    model_attributes = model_instance.public_methods(false).select { |m| m.to_s.end_with?('=') }.map { |m| m.to_s.chomp('=').to_sym }

    documented_attributes = properties.map { |p| p.to_s.underscore.to_sym }

    missing_from_swagger = model_attributes - documented_attributes - ignored_attributes

    expect(missing_from_swagger).to be_empty, "The following model attributes are not documented in .to_swagger_schema: #{missing_from_swagger.join(', ')}"
  end
end
