class Campaign < ApplicationRecord
  has_many :subscriptions
  has_many :emails
end
