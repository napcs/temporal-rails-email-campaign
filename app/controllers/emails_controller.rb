class EmailsController < ApplicationController
  before_action :set_campaign
  before_action :set_email, only: %i[ show edit update destroy ]

  # GET /emails or /emails.json
  def index
    @emails = Email.all
  end

  # GET /emails/1 or /emails/1.json
  def show
  end

  # GET /emails/new
  def new
    @email = Email.new
  end

  # GET /emails/1/edit
  def edit
  end

  # POST /emails or /emails.json
  def create
    @email = Email.new(email_params)

    respond_to do |format|
      if @email.save
        format.html { redirect_to campaign_email_url(@campaign, @email), notice: "Email was successfully created." }
        format.json { render :show, status: :created, location: [@campaign, @email] }
      else
        format.html { render :new, status: :unprocessable_entity }
        format.json { render json: @email.errors, status: :unprocessable_entity }
      end
    end
  end

  # PATCH/PUT /emails/1 or /emails/1.json
  def update
    respond_to do |format|
      if @email.update(email_params)
        format.html { redirect_to campaign_email_url(@campaign, @email), notice: "Email was successfully updated." }
        format.json { render :show, status: :ok, location: [@campaign, @email] }
      else
        format.html { render :edit, status: :unprocessable_entity }
        format.json { render json: @email.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /emails/1 or /emails/1.json
  def destroy
    @email.destroy

    respond_to do |format|
      format.html { redirect_to campaign_emails_url(@campaign), notice: "Email was successfully destroyed." }
      format.json { head :no_content }
    end
  end

  private

  # Use callbacks to share common setup or constraints between actions.
  def set_email
    @email = Email.find(params[:id])
  end

  # Only allow a list of trusted parameters through.
  def email_params
    params.require(:email).permit(:name, :body, :order)
  end

  def set_campaign
    @campaign = Campaign.find(params[:campaign_id])
  end
end
