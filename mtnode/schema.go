// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package mtnode

var (
	blockValidatorPrefix     string = "v"         // the prefix for all block validator keys
	messagePrefix            []byte = []byte("m") // maps a message sequence number to a message
	delayedMessagePrefix     []byte = []byte("d") // maps a delayed sequence number to an accumulator and a message
	sequencerBatchMetaPrefix []byte = []byte("s") // maps a batch sequence number to BatchMetadata
	delayedSequencedPrefix   []byte = []byte("a") // maps a delayed message count to the first sequencer batch sequence number with this delayed count

	messageCountKey        []byte = []byte("_messageCount")        // contains the current message count
	delayedMessageCountKey []byte = []byte("_delayedMessageCount") // contains the current delayed message count
	sequencerBatchCountKey []byte = []byte("_sequencerBatchCount") // contains the current sequencer message count
)
