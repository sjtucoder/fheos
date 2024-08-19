// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE, euint8, euint16, euint32, euint64, euint128, euint256, ebool} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract LteTest {
    using Utils for *;
    
    function lte(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "lte(euint8,euint8)")) {
            if (FHE.decrypt(FHE.lte(FHE.asEuint8(a), FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "lte(euint16,euint16)")) {
            if (FHE.decrypt(FHE.lte(FHE.asEuint16(a), FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "lte(euint32,euint32)")) {
            if (FHE.decrypt(FHE.lte(FHE.asEuint32(a), FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "lte(euint64,euint64)")) {
            if (FHE.decrypt(FHE.lte(FHE.asEuint64(a), FHE.asEuint64(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "lte(euint128,euint128)")) {
            if (FHE.decrypt(FHE.lte(FHE.asEuint128(a), FHE.asEuint128(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint8.lte(euint8)")) {
            if (FHE.asEuint8(a).lte(FHE.asEuint8(b)).decrypt()) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint16.lte(euint16)")) {
            if (FHE.asEuint16(a).lte(FHE.asEuint16(b)).decrypt()) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint32.lte(euint32)")) {
            if (FHE.asEuint32(a).lte(FHE.asEuint32(b)).decrypt()) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint64.lte(euint64)")) {
            if (FHE.asEuint64(a).lte(FHE.asEuint64(b)).decrypt()) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "euint128.lte(euint128)")) {
            if (FHE.asEuint128(a).lte(FHE.asEuint128(b)).decrypt()) {
                return 1;
            }
            return 0;
        }
        revert TestNotFound(test);
    }
}
