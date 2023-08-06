// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;


/**
 * @title L1Burn
 * @notice L1Burn keeps track of the total amount of ETH burned on L1.
 */
contract L1Burn {
    /**
     * @notice Total amount of ETH burned on L1.
     */
    uint256 public total;

    /**
     * @notice Mapping of blocks numbers to total burn.
     */
    mapping (uint64 => uint256) public reports;

    /**
     * @notice Allows the system address to submit a report.
     *
     * @param _blocknum L1 block number the report corresponds to.
     * @param _burn     Amount of ETH burned in the block.
     */
    function report(uint64 _blocknum, uint64 _burn) public {
        require(
            msg.sender == 0xDeaDDEaDDeAdDeAdDEAdDEaddeAddEAdDEAd0001,
            "L1Burn: reports can only be made from system address"
        );

        total += _burn;
        reports[_blocknum] = total;
    }

    /**
     * @notice Tallies up the total burn since a given block number.
     *
     * @param _blocknum L1 block number to tally from.
     *
     * @return Total amount of ETH burned since the given block number;
     */
    function tally(uint64 _blocknum) public view returns (uint256) {
        return total - reports[_blocknum];
    }
}
