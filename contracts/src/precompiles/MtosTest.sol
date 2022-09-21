// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE
// SPDX-License-Identifier: BUSL-1.1

pragma solidity >=0.4.21 <0.9.0;

/// @title Deprecated - Provides a method of burning arbitrary amounts of gas,
/// @notice This exists for historical reasons. Pre-Mantle, `MtosTest` had additional methods only the zero address could call.
/// These have been removed since users don't use them and calls to missing methods revert.
/// Precompiled contract that exists in every Mantle chain at 0x0000000000000000000000000000000000000069.
interface MtosTest {
    /// @notice Unproductively burns the amount of L2 ArbGas
    function burnArbGas(uint256 gasAmount) external pure;
}
