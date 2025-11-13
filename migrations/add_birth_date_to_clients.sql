-- Migration: Adicionar campo birth_date à tabela clients
-- Data: 2025-11-13

-- Adicionar coluna birth_date (permitir NULL temporariamente)
ALTER TABLE clients 
ADD COLUMN birth_date TIMESTAMP;

-- Atualizar registros existentes com data padrão
UPDATE clients 
SET birth_date = '1990-01-01 00:00:00'
WHERE birth_date IS NULL;

-- Tornar a coluna NOT NULL
ALTER TABLE clients 
ALTER COLUMN birth_date SET NOT NULL;

-- Comentário
COMMENT ON COLUMN clients.birth_date IS 'Data de nascimento do cliente';

