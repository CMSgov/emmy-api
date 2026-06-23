require 'rails_helper'

module Education
  RSpec.describe EnrollmentRequestV0, type: :model do
    it_behaves_like 'swagger_schema_drift_detection', EnrollmentRequestV0
  end
end
