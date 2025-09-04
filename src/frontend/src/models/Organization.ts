export interface OrganizationRequest {
  name: string;
  address: string;
  vat_number: string;
  state: string;
  city: string;
  zip_code: string;
  contact_email: string;
  pec_email: string;
  sdi_code: string;
  contact_phone: string;
  personnel_info: string;
  user_id: string;
}

export interface OrganizationResponse {
  message: string;
  organization_id: string;
}
