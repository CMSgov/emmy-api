require 'swagger_helper'

RSpec.describe 'api/v0/veteran_disability_ratings', type: :request do
  path '/api/v0/veteran-disability-ratings' do
    before do
      fake_result = Veteran::DisabilityRatingResponse.new(
        totalDisabilityStatus: true,
        totalDisabilityStatusEffectiveDate: "2023-01-01",
        combinedDisabilityRating: 100,
        combinedEffectiveDate: "2023-01-01",
        legalEffectiveDate: "2023-01-01",
        earliestRatingEndDate: "2024-06-01",
        dataSource: "VA",
        metadata: {
          durationMs: 125
        }
      )
      coordinator_mock = instance_double(Veteran::ServiceCoordinator)
      expect(coordinator_mock).to receive(:lookup_disability_rating)
                                    .with(
                                      {
                                        firstName: "John",
                                        lastName: "Doe",
                                        dateOfBirth: "1990-01-01",
                                        ssn: "000-00-0000"
                                      }
                                    )
                                    .and_return(fake_result)

      allow(Veteran::ServiceCoordinator).to receive(:new).and_return(coordinator_mock)
    end

    post('create veteran-disability-rating') do
      security [ oauth2: [] ]
      let(:Authorization) { 'Bearer <token>' }
      consumes 'application/json'
      produces 'application/json'
      parameter name: :veteran_disability_rating,
                in: :body,
                schema: Veteran::DisabilityRatingRequest.to_swagger_schema

      response(200, 'successful') do
        schema Veteran::DisabilityRatingResponse.to_swagger_schema

        let(:veteran_disability_rating) do
          {
            firstName: 'John',
            lastName: 'Doe',
            dateOfBirth: '1990-01-01',
            ssn: '000-00-0000'
          }
        end

        example 'application/json', :successful_disability_rating_lookup, {
          totalDisabilityStatus: true,
          totalDisabilityStatusEffectiveDate: "2023-01-01",
          combinedDisabilityRating: 100,
          combinedEffectiveDate: "2023-01-01",
          legalEffectiveDate: "2023-01-01",
          earliestRatingEndDate: "2024-06-01",
          dataSource: "VA",
          metadata: {
            durationMs: 125
          }
        }

        after do |example|
          example.metadata[:response][:content] = {
            'application/json' => {
              example: JSON.parse(response.body, symbolize_names: true)
            }
          }
        end
        run_test!
      end
    end
  end
end
