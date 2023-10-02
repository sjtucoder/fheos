pragma solidity >=0.4.21 <0.9.0;

interface FheOps {
    function lior(uint32 a, uint32 b) external view returns (uint32);

    function moshe(
        bytes memory input,
        uint32 inputLen
    ) external view returns (bytes32[1] memory);
}
