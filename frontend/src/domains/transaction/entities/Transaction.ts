import type { Ticket } from '../../ticket/entities/Ticket';
import type { User } from '../../auth/entities/User';

export interface Transaction {
  id: number;
  ticket_id: number;
  buyer_id: number;
  seller_id: number;
  status: TransactionStatus;
  ticket?: Ticket;
  buyer?: User;
  seller?: User;
}

export type TransactionStatus = 'Pending' | 'Waiting_Payment' | 'Verifying' | 'Completed';
