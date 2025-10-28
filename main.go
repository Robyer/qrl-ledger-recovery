package main

import (
	"bufio"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	qrlcommon "github.com/theQRL/go-qrllib/common"
	xmss "github.com/theQRL/go-qrllib/crypto/xmss"
	xmsswallet "github.com/theQRL/go-qrllib/legacywallet/xmss"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║               QRL Ledger Recovery Tool                     ║")
	fmt.Println("║     Extract QRL Keys from BIP39 Mnemonic (Offline)         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Purpose: Recover QRL from Ledger Nano S devices that no longer")
	fmt.Println("         support the QRL app due to firmware v2.x incompatibility.")
	fmt.Println()

	// Display security warnings
	fmt.Println(strings.Repeat("═", 60))
	fmt.Println("⚠️  CRITICAL SECURITY WARNINGS - READ CAREFULLY")
	fmt.Println(strings.Repeat("═", 60))
	fmt.Println()
	fmt.Println("🔴 EXPOSURE OF ALL ASSETS:")
	fmt.Println("   Entering your BIP39 mnemonic here will compromise the security")
	fmt.Println("   of ALL cryptocurrencies on your Ledger device, not just QRL!")
	fmt.Println("   You should move all other assets to new wallets after using")
	fmt.Println("   this tool.")
	fmt.Println()
	fmt.Println("🔴 COMPUTER SECURITY:")
	fmt.Println("   • NEVER use this tool on a public, shared, or compromised computer")
	fmt.Println("   • DISCONNECT from the internet before entering your mnemonic")
	fmt.Println("   • This tool works completely offline for your protection")
	fmt.Println("   • Verify your computer is free of malware/keyloggers")
	fmt.Println()
	fmt.Println("🔴 MNEMONIC SECURITY:")
	fmt.Println("   • Your BIP39 mnemonic gives FULL ACCESS to ALL funds")
	fmt.Println("   • NEVER share your mnemonic or hexseed with anyone")
	fmt.Println("   • NEVER enter your mnemonic on websites or unknown software")
	fmt.Println("   • Under normal circumstances, mnemonics should ONLY be entered")
	fmt.Println("     into hardware wallets - this is an emergency recovery tool")
	fmt.Println()
	fmt.Println("✅ VERIFICATION REQUIRED:")
	fmt.Println("   • Verify the QRL addresses shown match your Ledger device")
	fmt.Println("   • If using a passphrase on your Ledger, enter it correctly here")
	fmt.Println("   • Double-check all information before moving funds")
	fmt.Println()
	fmt.Println("📋 RECOMMENDED STEPS:")
	fmt.Println("   1. Disconnect your computer from the internet for maximum security")
	fmt.Println("   2. Run this tool to extract your QRL keys")
	fmt.Println("   3. Import QRL keys into official QRL wallet")
	fmt.Println("   4. Transfer QRL to a new secure wallet")
	fmt.Println("   5. Move all other cryptocurrencies to new wallets with new seeds")
	fmt.Println("   6. Consider your old Ledger seed fully compromised")
	fmt.Println()
	fmt.Println(strings.Repeat("═", 60))

	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Print("Do you understand these risks and accept full responsibility? (yes/no): ")
	confirmInput, _ := reader.ReadString('\n')
	confirm := strings.ToLower(strings.TrimSpace(confirmInput))
	if confirm != "yes" && confirm != "y" {
		return
	}
	fmt.Println()

	// Read BIP39 mnemonic
	fmt.Println(strings.Repeat("─", 60))
	fmt.Print("Enter your BIP39 mnemonic (12, 15, 18, 21, or 24 words): ")
	mnemonicInput, _ := reader.ReadString('\n')
	mnemonic := strings.TrimSpace(mnemonicInput)

	// Validate BIP39 mnemonic
	if !bip39.IsMnemonicValid(mnemonic) {
		fmt.Println("\n❌ Error: Invalid BIP39 mnemonic")
		fmt.Println("\nPlease verify your mnemonic and try again.")
		fmt.Println("Ensure all words are spelled correctly and in the right order.")
		return
	}

	wordCount := len(strings.Fields(mnemonic))
	fmt.Printf("\n✅ Valid BIP39 mnemonic detected (%d words)\n", wordCount)

	// Ask for optional passphrase
	fmt.Println()
	fmt.Println("BIP39 Passphrase (25th word):")
	fmt.Println("If you set a passphrase on your Ledger, enter it exactly as configured.")
	fmt.Print("Otherwise, press Enter to continue with no passphrase: ")
	passphraseInput, _ := reader.ReadString('\n')
	passphrase := strings.TrimRight(passphraseInput, "\r\n")

	if passphrase == "" {
		fmt.Println("\n✅ No passphrase - using standard derivation")
	} else {
		fmt.Printf("\n✅ Passphrase entered (%d characters)\n", len(passphrase))
		fmt.Println("⚠️  Ensure this matches your Ledger passphrase exactly!")
	}

	// Derive keys for both trees
	fmt.Println()
	fmt.Println("Deriving QRL keys from BIP39 seed...")
	tree1Key, tree2Key, err := deriveTreeKeys(mnemonic, passphrase)
	if err != nil {
		fmt.Printf("\n❌ Error deriving tree keys: %v\n", err)
		return
	}

	// Derive QRL keys
	fmt.Println("🌲 Generating Tree 1 (this may take a moment)...")
	qrl1Key, err := deriveQRLKey(tree1Key)
	if err != nil {
		fmt.Printf("\n❌ Error deriving QRL key for Tree 1: %v\n", err)
		return
	}

	fmt.Println("🌲 Generating Tree 2 (this may take a moment)...")
	qrl2Key, err := deriveQRLKey(tree2Key)
	if err != nil {
		fmt.Printf("\n❌ Error deriving QRL key for Tree 2: %v\n", err)
		return
	}

	// Pre-display confirmation
	fmt.Println()
	fmt.Println("✅ QRL keys successfully generated!")
	fmt.Println()
	fmt.Print("Press Enter when ready to display keys...")
	_, _ = reader.ReadString('\n')

	fmt.Println()

	// Display results
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                  QRL RECOVERY KEYS                         ║")
	fmt.Println("║             !!! PRIVATE - KEEP SECURE !!!                  ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	displaySingleTree(qrl1Key, 1)
	displaySingleTree(qrl2Key, 2)

	fmt.Println()
	fmt.Println(strings.Repeat("═", 60))
	fmt.Println("🔐 NEXT STEPS:")
	fmt.Println(strings.Repeat("═", 60))
	fmt.Println("1. Verify the addresses above match those on your Ledger device")
	fmt.Println("2. Import these keys into the official QRL wallet")
	fmt.Println("3. Transfer your QRL to a new secure wallet")
	fmt.Println("4. Move all other cryptocurrencies to new wallets")
	fmt.Println("5. Consider your old Ledger seed permanently compromised")
	fmt.Println()
	fmt.Println("⚠️  Your old Ledger seed should NEVER be used again!")
	fmt.Println(strings.Repeat("═", 60))
	fmt.Println()

	fmt.Print("Press Enter to exit...")
	_, _ = reader.ReadString('\n')
}

func displaySingleTree(keys *QRLKeys, treeNumber int) {
	fmt.Println()
	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("🌲 TREE %d\n", treeNumber)
	fmt.Println(strings.Repeat("─", 60))

	fmt.Println("\n📍 QRL Address (PUBLIC):")
	fmt.Printf("   %s\n", keys.Address)

	fmt.Println("\n🔑 QRL Mnemonic (PRIVATE - 34 words):")
	fmt.Printf("   %s\n", keys.Mnemonic)

	fmt.Println("\n🔑 QRL Hexseed (PRIVATE - 102 hex characters):")
	fmt.Printf("   %s\n", keys.Hexseed)
}

type QRLKeys struct {
	Address  string
	Mnemonic string
	Hexseed  string
}

func deriveTreeKeys(mnemonic, passphrase string) (*bip32.Key, *bip32.Key, error) {
	// Generate BIP39 seed with passphrase
	seed := bip39.NewSeed(mnemonic, passphrase)

	// Create master key using BIP32
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create master key: %w", err)
	}

	// Derive path: m/44'/238'/0'/0'/[treeNumber]'

	// m/44'
	purposeKey, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44) // BIP44 Purpose (44' hardened)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive purpose: %w", err)
	}

	// m/44'/238'
	coinTypeKey, err := purposeKey.NewChildKey(bip32.FirstHardenedChild + 238) // QRL BIP44 coin type (238' hardened)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive coin type: %w", err)
	}

	// m/44'/238'/0'
	accountKey, err := coinTypeKey.NewChildKey(bip32.FirstHardenedChild + 0) // Account (0' hardened)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive account: %w", err)
	}

	// m/44'/238'/0'/0'
	changeKey, err := accountKey.NewChildKey(bip32.FirstHardenedChild + 0) // Change (0' hardened)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive change: %w", err)
	}

	// m/44'/238'/0'/0'/0' for Tree 1
	address1Key, err := changeKey.NewChildKey(bip32.FirstHardenedChild + 0) // Tree 1 address (0' hardened)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive address for tree 1: %w", err)
	}

	// m/44'/238'/0'/0'/1' for Tree 2
	address2Key, err := changeKey.NewChildKey(bip32.FirstHardenedChild + 1) // Tree 2 address (1' hardened)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive address for tree 2: %w", err)
	}

	return address1Key, address2Key, nil
}

func deriveQRLKey(treeKey *bip32.Key) (*QRLKeys, error) {
	// Concatenate private key (32 bytes) + chain code (32 bytes)
	combined := make([]byte, 64)
	copy(combined[0:32], treeKey.Key)        // Private key
	copy(combined[32:64], treeKey.ChainCode) // Chain code

	// Hash with SHA3-512
	hasher := sha3.New512()
	hasher.Write(combined)
	hashed := hasher.Sum(nil) // 64 bytes

	// Take first 48 bytes as QRL seed
	seed := [48]uint8(hashed[0:48])

	// Use specific parameters used by Ledger app
	height := xmss.ToHeight(8)
	hashFunc := xmss.SHA2_256
	addrType := qrlcommon.SHA256_2X

	// Create QRL XMSS wallet. This call can take a bit of time as it is initializing whole XMSS tree
	wallet := xmsswallet.NewWalletFromSeed(seed, height, hashFunc, addrType)

	// Return results
	addr := wallet.GetAddress()
	eSeed := wallet.GetExtendedSeed()

	return &QRLKeys{
		Hexseed:  hex.EncodeToString(eSeed[:]),
		Mnemonic: wallet.GetMnemonic(),
		Address:  "Q" + hex.EncodeToString(addr[:]),
	}, nil
}
