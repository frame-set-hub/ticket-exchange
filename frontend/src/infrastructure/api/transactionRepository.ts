import type { ITransactionRepository } from '../../domains/transaction/repositories/TransactionRepository';
import type { Transaction } from '../../domains/transaction/entities/Transaction';
import apiClient from './apiClient';

export class TransactionRepository implements ITransactionRepository {
  async initiate(ticketId: number): Promise<Transaction> {
    const { data } = await apiClient.post(`/transactions/${ticketId}`);
    return data;
  }

  async uploadPayment(txId: number): Promise<void> {
    await apiClient.post(`/transactions/${txId}/upload-payment`);
  }

  async uploadTicket(txId: number): Promise<void> {
    await apiClient.post(`/transactions/${txId}/upload-ticket`);
  }

  async getAll(): Promise<Transaction[]> {
    const { data } = await apiClient.get('/admin/transactions');
    return data;
  }

  async complete(txId: number): Promise<void> {
    await apiClient.post(`/admin/transactions/${txId}/complete`);
  }
}
