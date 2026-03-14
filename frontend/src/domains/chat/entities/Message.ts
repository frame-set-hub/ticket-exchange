export interface Message {
  id: number;
  transaction_id: number;
  sender_id: number;
  sender_username?: string;
  receiver_id: number;
  content: string;
  attachment_url?: string;
  created_at: string;
}

export interface SendMessageParams {
  content: string;
  attachment_url?: string;
}
