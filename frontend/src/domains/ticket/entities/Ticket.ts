export interface Ticket {
  id: number;
  title: string;
  venue: string;
  price: number;
  category: string;
  description?: string;
  status: TicketStatus;
  seller: {
    id: number;
    username: string;
  };
}

export type TicketStatus = 'Available' | 'Pending' | 'Waiting_Payment' | 'Verifying' | 'Completed';

export interface TicketFilters {
  title?: string;
  category?: string;
  min_price?: string;
  max_price?: string;
}

export interface CreateTicketParams {
  title: string;
  venue: string;
  price: number;
  category: string;
  description: string;
}
