json.extract! subscription, :id, :email, :campaign_id, :status, :created_at, :updated_at
json.url subscription_url(subscription, format: :json)
