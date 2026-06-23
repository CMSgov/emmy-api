require 'rails_helper'

module Education
  RSpec.describe EnrollmentResponseV0, type: :model do
    it_behaves_like 'swagger_schema_drift_detection', EnrollmentResponseV0, ignored_attributes: [ :raw_data ]
  end
end
