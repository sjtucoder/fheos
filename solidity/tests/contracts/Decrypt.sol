// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE, euint8, euint16, euint32, euint64, euint128, euint256, ebool} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract DecryptTest {
    using Utils for *;
    
    function decrypt(string calldata test) public {
        if (Utils.cmp(test, "decrypt(euint8)")) {
            euint8 aEnc = FHE.asEuint8(1);
            require(aEnc.decrypt(1) < 10, "Decryption failed");
            require(FHE.decrypt(aEnc, 1) < 10, "Decryption failed");
        } else if (Utils.cmp(test, "decrypt(euint16)")) {
            euint16 aEnc = FHE.asEuint16(1);
            require(aEnc.decrypt(1) < 10, "Decryption failed");
            require(FHE.decrypt(aEnc, 1) < 10, "Decryption failed");
        } else if (Utils.cmp(test, "decrypt(euint32)")) {
            euint32 aEnc = FHE.asEuint32(1);
            require(aEnc.decrypt(1) < 10, "Decryption failed");
            require(FHE.decrypt(aEnc, 1) < 10, "Decryption failed");
        } else if (Utils.cmp(test, "decrypt(euint64)")) {
            euint64 aEnc = FHE.asEuint64(1);
            require(aEnc.decrypt(1) < 10, "Decryption failed");
            require(FHE.decrypt(aEnc, 1) < 10, "Decryption failed");
        } else if (Utils.cmp(test, "decrypt(euint128)")) {
            euint128 aEnc = FHE.asEuint128(1);
            require(aEnc.decrypt(1) < 10, "Decryption failed");
            require(FHE.decrypt(aEnc, 1) < 10, "Decryption failed");
        } else if (Utils.cmp(test, "decrypt(euint256)")) {
            euint256 aEnc = FHE.asEuint256(1);
            require(aEnc.decrypt(1) < 10, "Decryption failed");
            require(FHE.decrypt(aEnc, 1) < 10, "Decryption failed");
        } else if (Utils.cmp(test, "decrypt(ebool)")) {
            ebool aEnc = FHE.asEbool(false);
            require(!aEnc.decrypt(false), "Decryption failed");
            require(!FHE.decrypt(aEnc, false), "Decryption failed");
        } else if (Utils.cmp(test, "decrypt(euint8) fail")) {
            euint8 aEnc = FHE.asEuint8(1);
            require(aEnc.decrypt() < 10, "Decryption failed");
            require(FHE.decrypt(aEnc) < 10, "Decryption failed");
        } else if (Utils.cmp(test, "decrypt(euint16) fail")) {
            euint16 aEnc = FHE.asEuint16(1);
            require(aEnc.decrypt() < 10, "Decryption failed");
            require(FHE.decrypt(aEnc) < 10, "Decryption failed");
        } else if (Utils.cmp(test, "decrypt(euint32) fail")) {
            euint32 aEnc = FHE.asEuint32(1);
            require(aEnc.decrypt() < 10, "Decryption failed");
            require(FHE.decrypt(aEnc) < 10, "Decryption failed");
        } else if (Utils.cmp(test, "decrypt(euint64) fail")) {
            euint64 aEnc = FHE.asEuint64(1);
            require(aEnc.decrypt() < 10, "Decryption failed");
            require(FHE.decrypt(aEnc) < 10, "Decryption failed");
        } else if (Utils.cmp(test, "decrypt(euint128) fail")) {
            euint128 aEnc = FHE.asEuint128(1);
            require(aEnc.decrypt() < 10, "Decryption failed");
            require(FHE.decrypt(aEnc) < 10, "Decryption failed");
        } else if (Utils.cmp(test, "decrypt(euint256) fail")) {
            euint256 aEnc = FHE.asEuint256(1);
            require(aEnc.decrypt() < 10, "Decryption failed");
            require(FHE.decrypt(aEnc) < 10, "Decryption failed");
        }  else {
            revert TestNotFound(test);
        }
    }
}
