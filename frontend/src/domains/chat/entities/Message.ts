export interface Message {
  transaction_id: number;
  sender_id: number;
  receiver_id: number;
  content: string;
}

export interface SendMessageParams {
  transaction_id: number;
  content: string;
  receiver_id: number;
}
