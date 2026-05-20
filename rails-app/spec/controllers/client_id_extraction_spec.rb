require 'rails_helper'

RSpec.describe ApplicationController, type: :controller do
  controller do
    def index
      render json: { client_id: client_id }
    end
  end

  describe '#client_id' do
    context 'when X-Sub header is present' do
      it 'returns the value from X-Sub' do
        request.headers['X-Sub'] = 'test-sub'
        get :index
        expect(JSON.parse(response.body)['client_id']).to eq('test-sub')
      end
    end

    context 'when Authorization header with Bearer token is present' do
      let(:payload) { { sub: 'jwt-user-123' } }
      let(:token) { JWT.encode(payload, nil, 'none') }

      it 'returns the sub claim from the JWT' do
        request.headers['Authorization'] = "Bearer #{token}"
        get :index
        expect(JSON.parse(response.body)['client_id']).to eq('jwt-user-123')
      end

      context 'when sub claim is missing' do
        let(:payload) { { other: 'claim' } }
        it 'returns unknown-subject' do
          request.headers['Authorization'] = "Bearer #{token}"
          get :index
          expect(JSON.parse(response.body)['client_id']).to eq('unknown-subject')
        end
      end

      context 'when token is malformed' do
        it 'returns unknown-subject' do
          request.headers['Authorization'] = "Bearer invalid.token.here"
          get :index
          expect(JSON.parse(response.body)['client_id']).to eq('unknown-subject')
        end
      end
    end

    context 'when no relevant headers are present' do
      it 'returns unknown-subject' do
        get :index
        expect(JSON.parse(response.body)['client_id']).to eq('unknown-subject')
      end
    end

    context 'when Authorization header is present but not Bearer' do
      it 'returns unknown-subject' do
        request.headers['Authorization'] = "Basic YWxhZGRpbjpvcGVuc2VzYW1l"
        get :index
        expect(JSON.parse(response.body)['client_id']).to eq('unknown-subject')
      end
    end
  end
end
