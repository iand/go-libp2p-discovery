package discovery

import (
	"context"
	"time"

	logging "github.com/ipfs/go-log"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

var log = logging.Logger("discovery")

// FindPeers is a utility function that synchonously collects peers from a Discoverer
func FindPeers(ctx context.Context, d Discoverer, ns string, opts ...Option) ([]pstore.PeerInfo, error) {
	var res []pstore.PeerInfo

	ch, err := d.FindPeers(ctx, ns, opts...)
	if err != nil {
		return nil, err
	}

	for pi := range ch {
		res = append(res, pi)
	}

	return res, nil
}

// Advertise is a utility function that persistently advertises a service through an Advertiser
func Advertise(ctx context.Context, a Advertiser, ns string, opts ...Option) {
	go func() {
		for {
			ttl, err := a.Advertise(ctx, ns, opts...)
			if err != nil {
				log.Debugf("Error advertising %s: %s", ns, err.Error())
				if ctx.Err() != nil {
					return
				}

				select {
				case <-time.After(2 * time.Minute):
					continue
				case <-ctx.Done():
					return
				}
			}

			wait := 7 * ttl / 8
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				return
			}
		}
	}()
}
