// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE
// SPDX-License-Identifier: BUSL-1.1

pragma solidity ^0.8.0;

interface IChallengeResultReceiver {
    function completeChallenge(
        uint256 challengeIndex,
        address winner,
        address loser
    ) external;
}
