import { useState, useEffect, useCallback } from 'react';
import type { Transaction } from '../../../domains/transaction/entities/Transaction';
import { getServices } from '../../../infrastructure/services/ServiceContainer';

export function useEscrow(ticketId: number) {
  const [tx, setTx] = useState<Transaction | null>(null);
  const { transactionRepository } = getServices();

  const initEscrow = useCallback(async () => {
    try {
      const data = await transactionRepository.initiate(ticketId);
      setTx(data);
    } catch {
      // Transaction may already exist
    }
  }, [ticketId]);

  const uploadProof = async (userId: number) => {
    if (!tx) return;
    try {
      if (userId === tx.buyer_id) {
        await transactionRepository.uploadPayment(tx.id);
      } else if (userId === tx.seller_id) {
        await transactionRepository.uploadTicket(tx.id);
      }
      alert('Uploaded successfully to Escrow!');
    } catch {
      console.error('Upload failed');
    }
  };

  useEffect(() => {
    initEscrow();
  }, [initEscrow]);

  return { tx, uploadProof };
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
    try {
      await transactionRepository.complete(txId);
      fetchTransactions();
    } catch {
      alert('Verification failed or not ready.');
    }
  };

  useEffect(() => {
    fetchTransactions();
  }, [fetchTransactions]);

  return { transactions, completeTransaction };
}
