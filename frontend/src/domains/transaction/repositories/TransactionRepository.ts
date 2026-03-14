import type { Transaction } from '../entities/Transaction';

export interface ITransactionRepository {
  initiate(ticketId: number): Promise<Transaction>;
  getByTicketId(ticketId: number): Promise<Transaction>;
  uploadPayment(txId: number): Promise<void>;
  uploadTicket(txId: number): Promise<void>;
  getAll(): Promise<Transaction[]>;
  complete(txId: number): Promise<void>;
}
