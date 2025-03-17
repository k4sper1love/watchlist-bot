// Package security provides utilities for secure encryption and decryption of sensitive data.
//
// It uses the AES-GCM (Galois/Counter Mode) encryption algorithm, which ensures both confidentiality
// and integrity of the data. The package is designed to be simple, secure, and easy to integrate
// into applications that require encryption capabilities.
//
// Features:
//   - AES-GCM Encryption: Uses AES-256 in GCM mode for authenticated encryption.
//     This ensures data cannot be tampered with without detection.
//   - Master Key Management: Retrieves the master key from an environment variable ("MASTER_KEY"),
//     which must be exactly 32 bytes long for AES-256 compatibility.
//   - Nonce Generation: Creates a unique nonce (number used once) for each encryption operation,
//     ensuring security against replay attacks.
//   - Base64 Encoding: Encodes encrypted data in base64 for safe storage and transmission.
//   - Error Handling: Provides detailed error messages for invalid keys, malformed inputs, and failures.
//
// Usage:
//
// First, ensure the `MASTER_KEY` environment variable is set to a 32-byte key.
// Then, use `Encrypt` to encrypt sensitive data and `Decrypt` to recover it.
//
// Example:
//
//	encrypted, err := security.Encrypt("sensitive data")
//	if err != nil {
//	    log.Fatalf("Encryption failed: %v", err)
//	}
//	fmt.Println("Encrypted:", encrypted)
//
//	decrypted, err := security.Decrypt(encrypted)
//	if err != nil {
//	    log.Fatalf("Decryption failed: %v", err)
//	}
//	fmt.Println("Decrypted:", decrypted)
//
// Security Considerations:
// - Key Management: Keep `MASTER_KEY` secret. Do not expose it in source code or logs.
// - Nonce Reuse: Each encryption operation requires a unique nonce. Reusing nonces compromises security.
// - Data Integrity: AES-GCM provides built-in integrity checks. If tampered with, decryption will fail.
//
// This package is suitable for securely handling sensitive data such as passwords, tokens, and confidential information.
package security
