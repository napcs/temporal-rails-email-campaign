require 'temporalio'

module Temporal
  def self.connection
    @connection ||= Temporalio::Connection.new('localhost:7233')
  end

  def self.task_queue_name
    'email_drips'
  end
end
