Rails.application.routes.draw do
  resources :subscriptions

  resources :campaigns do
    resources :emails
    member do
      get :welcome_message, defaults: { format: 'json' }
      get :unsubscribe_message, defaults: { format: 'json' }
    end
  end
  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html

  # Defines the root path route ("/")
  # root "articles#index"
end
