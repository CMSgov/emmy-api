require 'rails_helper'

module Education
  RSpec.describe EnrollmentResponse, type: :model do
    it_behaves_like 'swagger_schema_drift_detection', EnrollmentResponse, ignored_attributes: [ :raw_data ]
  end
end
