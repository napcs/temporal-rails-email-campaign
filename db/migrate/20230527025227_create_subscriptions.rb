class CreateSubscriptions < ActiveRecord::Migration[7.0]
  def change
    create_table :subscriptions do |t|
      t.string :email
      t.references :campaign, null: false, foreign_key: true
      t.string :status

      t.timestamps
    end
  end
end
