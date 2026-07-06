require 'swagger_helper'

unless defined?(Veteran::NotFoundError)
  module Veteran
    class NotFoundError < StandardError; end
  end
end

RSpec.describe 'api/v1/veteran_disability_ratings', type: :request do
  path '/api/v1/veteran-disability-ratings' do
    post('create veteran-disability-rating v1') do
      security [ oauth2: [] ]
      let(:Authorization) { 'Bearer <token>' }
      consumes 'application/json'
      produces 'application/json'
      parameter name: :veteran_disability_rating,
                in: :body,
                schema: Veteran::DisabilityRatingRequestV1.to_swagger_schema

      response(200, 'successful') do
        schema Veteran::DisabilityRatingResponseV1.to_swagger_schema

        let(:veteran_disability_rating) do
          {
            vadrRequest: {
              personGivenName: 'John',
              personSurName: 'Doe',
              personBirthDate: '1990-01-01',
              personSocialSecurityNumber: '000-00-0000',
              personSexCode: 'M',
              personContactInformation: {
                streetLineOneAddress: '123 Main St',
                cityName: 'Arlington',
                stateText: 'VA',
                zipCode: '22202',
                countryText: 'USA'
              }
            }
          }
        end

        before do
          fake_result = Veteran::DisabilityRatingResponseV1.new(
            totalDisabilityStatus: true,
            permanentDisabilityStatus: true,
            totalDisabilityStatusEffectiveDate: "2023-01-01",
            combinedDisabilityRating: 100,
            combinedEffectiveDate: "2023-01-01",
            legalEffectiveDate: "2023-01-01",
            earliestRatingEndDate: "2024-06-01",
            pensionAwardStatusIndicator: false,
            individualRatings: [
              {
                decisionText: "Service Connection",
                ratingEffectiveDate: "2023-01-01",
                ratingEndDate: "2024-01-01",
                ratingPercentage: 50,
                disabilityRatingId: "12345",
                staticIndicator: true
              }
            ]
          )
          coordinator_mock = instance_double(Veteran::ServiceCoordinator)
          expect(coordinator_mock).to receive(:lookup_disability_rating)
                                        .with(kind_of(Veteran::DisabilityRatingRequestV1), "V1")
                                        .and_return(fake_result)

          allow(Veteran::ServiceCoordinator).to receive(:new).and_return(coordinator_mock)
        end

        run_test!
      end

      response(400, 'bad request - missing fields') do
        let(:veteran_disability_rating) do
          {
            vadrRequest: {
              personGivenName: 'John'
              # missing surName and birthDate
            }
          }
        end

        run_test! do |response|
          data = JSON.parse(response.body)
          expect(data['error']).to include('missing required field')
        end
      end

      response(404, 'not found') do
        let(:veteran_disability_rating) do
          {
            vadrRequest: {
              personGivenName: 'Unknown',
              personSurName: 'Veteran',
              personBirthDate: '1900-01-01'
            }
          }
        end

        before do
          coordinator_mock = instance_double(Veteran::ServiceCoordinator)
          allow(coordinator_mock).to receive(:lookup_disability_rating)
                                       .and_raise(Veteran::NotFoundError.new("Veteran not found"))
          allow(Veteran::ServiceCoordinator).to receive(:new).and_return(coordinator_mock)
        end

        run_test!
      end
    end
  end
end
