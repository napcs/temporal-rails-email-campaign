json.extract! email, :id, :name, :body, :order, :created_at, :updated_at
json.url campaign_email_url(@campaign, email, format: :json)
