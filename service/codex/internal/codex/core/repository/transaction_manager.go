package repository

type TransactionManager struct {
	CodexRepo              CodexRepository
	CodexDocumentRepo      CodexDocumentRepository
	CodexDocumentChunkRepo CodexDocumentChunkRepository
	QARepo                 QARepository
}
