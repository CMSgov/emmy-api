require 'rails_helper'

module Veteran
  RSpec.describe DisabilityRatingMapper do
    let(:start_time) { Time.now }
    let(:datasource_duration) { 1500 }

    let(:total_disability_response) do
      {
        'data' => {
          'total_disability' => {
            'status' => true,
            'effective_date' => '2023-01-01'
          },
          'permanent_and_total' => {
            'service_connected_status' => true
          }
        }
      }
    end

    let(:disability_score_response) do
      {
        'data' => {
          'attributes' => {
            'combined_disability_rating' => 100,
            'combined_effective_date' => '2023-01-01',
            'legal_effective_date' => '2023-01-01',
            'individual_ratings' => [
              { 'rating_end_date' => '2025-01-01' },
              { 'rating_end_date' => '2024-06-01' },
              { 'rating_end_date' => '2026-01-01' }
            ]
          }
        }
      }
    end

    describe '.map_response' do
      subject { described_class.map_response(total_disability_response, disability_score_response, datasource_duration, start_time) }

      it 'returns a Veteran::DisabilityRatingResponse' do
        expect(subject).to be_a(Veteran::DisabilityRatingResponse)
      end

      it 'maps total disability fields correctly' do
        expect(subject.total_disability_status).to be true
        expect(subject.total_disability_status_effective_date).to eq('2023-01-01')
        expect(subject.permanent_disability_status).to be true
      end

      it 'maps combined disability fields correctly' do
        expect(subject.combined_disability_rating).to eq(100)
        expect(subject.combined_effective_date).to eq('2023-01-01')
        expect(subject.legal_effective_date).to eq('2023-01-01')
      end

      context 'when determining earliest_rating_end_date' do
        it 'picks the minimum date from individual_ratings' do
          expect(subject.earliest_rating_end_date).to eq('2024-06-01')
        end

        context 'when individual_ratings is empty' do
          let(:disability_score_response) do
            {
              'data' => {
                'attributes' => {
                  'individual_ratings' => []
                }
              }
            }
          end

          it 'returns nil' do
            expect(subject.earliest_rating_end_date).to be_nil
          end
        end

        context 'when individual_ratings contains blank or nil dates' do
          let(:disability_score_response) do
            {
              'data' => {
                'attributes' => {
                  'individual_ratings' => [
                    { 'rating_end_date' => '' },
                    { 'rating_end_date' => nil },
                    { 'rating_end_date' => '2025-05-01' }
                  ]
                }
              }
            }
          end

          it 'ignores blank/nil and picks the minimum valid date' do
            expect(subject.earliest_rating_end_date).to eq('2025-05-01')
          end
        end

        context 'when all rating_end_date values are blank' do
          let(:disability_score_response) do
            {
              'data' => {
                'attributes' => {
                  'individual_ratings' => [
                    { 'rating_end_date' => '' },
                    { 'rating_end_date' => nil }
                  ]
                }
              }
            }
          end

          it 'returns nil' do
            expect(subject.earliest_rating_end_date).to be_nil
          end
        end
      end

      it 'populates metadata correctly' do
        expect(subject.metadata[:requestTimestamp]).to eq(start_time.utc.iso8601(3))
        expect(subject.metadata[:datasourceDurationMillis]).to eq(1500)
        expect(subject.metadata[:transactionId]).to be_present
      end

      it 'includes raw data' do
        expect(subject.raw_data[:total_disability_response]).to eq(total_disability_response)
        expect(subject.raw_data[:disability_score_response]).to eq(disability_score_response)
      end
    end
  end
end
