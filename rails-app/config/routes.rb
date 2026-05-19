Rails.application.routes.draw do
  mount Rswag::Ui::Engine => "/api-docs"
  mount Rswag::Api::Engine => "/api-docs"
  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html

  # Reveal health status on /up that returns 200 if the app boots with no exceptions, otherwise 500.
  # Can be used by load balancers and uptime monitors to verify that the app is live.
  get "up" => "rails/health#show", as: :rails_health_check
  get "health" => "health#show"

  namespace :api do
    namespace :v0 do
      post "education-enrollments", to: "education_enrollments#create"
      post "batch-education-enrollments", to: "batch_education_enrollments#create"
      get "batch-education-enrollments/:id", to: "batch_education_enrollments#show"
      get "batch-education-enrollments/:batchJobId/details", to: "batch_education_enrollments#details"
      post "veteran-disability-ratings", to: "veteran_disability_ratings#create"
    end
  end

  # Defines the root path route ("/")
  # root "posts#index"
end
