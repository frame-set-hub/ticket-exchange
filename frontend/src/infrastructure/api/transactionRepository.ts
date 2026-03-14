import type { ITransactionRepository } from '../../domains/transaction/repositories/TransactionRepository';
import type { Transaction } from '../../domains/transaction/entities/Transaction';
import apiClient from './apiClient';

export class TransactionRepository implements ITransactionRepository {
  async initiate(ticketId: number): Promise<Transaction> {
    const { data } = await apiClient.post('/transactions', { ticket_id: ticketId });
    return data;
  }

  async getByTicketId(ticketId: number): Promise<Transaction> {
    const { data } = await apiClient.get(`/transactions/by-ticket/${ticketId}`);
    return data;
  }

  async uploadPayment(txId: number): Promise<void> {
    await apiClient.post(`/transactions/${txId}/status`, { status: 'Waiting_Payment' });
  }

  async uploadTicket(txId: number): Promise<void> {
    await apiClient.post(`/transactions/${txId}/status`, { status: 'Verifying' });
  }

  async getAll(): Promise<Transaction[]> {
    const { data } = await apiClient.get('/admin/transactions');
    return data;
  }

  async complete(txId: number): Promise<void> {
    await apiClient.post(`/admin/transactions/${txId}/status`, { status: 'Completed' });
  }
}
