import { useState, useEffect, useCallback } from 'react';
import type { Transaction } from '../../../domains/transaction/entities/Transaction';
import { getServices } from '../../../infrastructure/services/ServiceContainer';

export function useEscrow(ticketId: number) {
  const [tx, setTx] = useState<Transaction | null>(null);
  const [error, setError] = useState<string | null>(null);
  const { transactionRepository } = getServices();

  const initEscrow = useCallback(async () => {
    try {
      setError(null);
      // Try to find existing transaction for this ticket first
      const existing = await transactionRepository.getByTicketId(ticketId);
      setTx(existing);
    } catch {
      // No existing transaction — try to create one (buyer flow)
      try {
        const data = await transactionRepository.initiate(ticketId);
        setTx(data);
      } catch (err: any) {
        const msg = err.response?.data?.error || 'Failed to create transaction';
        setError(msg);
      }
    }
  }, [ticketId]);

  const uploadProof = async (userId: number) => {
    if (!tx) return;
    if (userId === tx.buyer_id) {
      await transactionRepository.uploadPayment(tx.id);
    } else if (userId === tx.seller_id) {
      await transactionRepository.uploadTicket(tx.id);
    }
  };

  useEffect(() => {
    initEscrow();
  }, [initEscrow]);

  return { tx, error, uploadProof };
}

export function useAdminTransactions() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const { transactionRepository } = getServices();

  const fetchTransactions = useCallback(async () => {
    try {
      const data = await transactionRepository.getAll();
      setTransactions(data);
    } catch (err) {
      console.error(err);
    }
  }, []);

  const completeTransaction = async (txId: number) => {
    await transactionRepository.complete(txId);
    fetchTransactions();
  };

  useEffect(() => {
    fetchTransactions();
  }, [fetchTransactions]);

  return { transactions, completeTransaction };
}
