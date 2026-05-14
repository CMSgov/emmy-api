require 'rails_helper'

module Education
  RSpec.describe EnrollmentRequest, type: :model do
    it_behaves_like 'swagger_schema_drift_detection', EnrollmentRequest
  end
end
