
type validator struct {
	genesisTime uint64 
	ticker *slotutil.SlotTicker 
	assigments *pbAssigmentsResponse
	proposerClient pb.ValidatorServiceClient
	validatorClient pb.ValidatorServiceClient 
	attersterClient ethpb.NodeClient
	node map[[48]byte]*keystore.Key 
	pubkeys [][]byte 
	prevBalance map[[48]byte]uint64
	logValidatorBalances bool
}

func (v *validator) Done()  {
	v.ticker.Done()
	
}

func (v *validaotr) WaitForChainStart(ctx context.Context) error  {
	ctx, span := trace.StartSpan(ctx, "validator.WaitFroChainStart")
	defer span.End()

	stream, err := v.validatorClient.WaitForChainStart(ctx, &ptypes.Empty{})
	if err != nil {
		return errors.Wrap(err, "coulde not")
	}
	for {
		log.Info()
		chainStartRes, err := streamRecv()

		if err == io.EOF {
			break
		}

		if ctx.Err() == context.Canceled {
			return errors.Wrap(ctx.Err(), "context has bben ")
		}
		if err != nil {
			return errors.Wrap(err, "")
		}
		v.genesisTime = chainStartRes.GenesisTime
		break
	}
	v.ticker = slotutil.GetSlotTicker(time.Unix(int64(v.genesisTime), 0), params.BeaconConfig().SecondPerSlot)
	log.WithField("genesisTime", time.Unix(uint64(v.genesisTime), 0)).Info("Beacon")
	return nil 
}


func (v *validator) WaitForActivation(ctx context.Context) error  {
	ctx, span := trace.StartSpan(ctx, "validator.WaitForActivation")
	defer span.End()
	req := &pb.ValidatorActivatationRequest{
		PublicKeys: v.pubkeys,
	}
	stream, err := v.validatorClient.WaitForActivation(ctx, req)
	if err != nil {
		return errors.Wrap(err, )
	}
	var ValidatorActivatatedRecords [][]byte 
	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if ctx.Err() == context.Canceled {
			return errors.Wrap(ctx.Err(), "")
		}
		if err != nil {
			return errors.Wrap(err, )
		}
		log.Info("Waiting")
		activatedkeys := v.checkAndLogActivationStatus(res.Statuses)

		if len(activatedKeys) > 0 {
			ValidatorActivatatedRecords = activatedKeys
			break 

		}
	}
	for _, pubkey := range ValidatorActivatatedRecords {
		log.WithField("pubkey", fmt.Sprintf)

	}
	v.ticker = slotutil.GetSlotTicker(time.Unix(int64(v.genesisTime), 0), params.BeaconConfig().SecondPerSlot)

	return nil 
}
func (v *validator) WaitForSync(ctx context.Context) error  {
	ctx, span := trace.StartSpan(ctx, "validator.WaitFor0")
	defer span.End()

	s, err := v.node.GetSyncStatus(ctx, &ptypes.Empty{})
	if err != nil {
		return errors.Wrap(err,)
	}
	if !s.Syncing {
		return nil
	}

	for {
		select {
		case <-time.After(time.Duration(params.BeaconConfig().SlotPerEpoch/2) * time.Second):
			S, err := v.node.GetSyncStatus(ctx, &ptypes.Empty{})
			if err != nil {
				return errors.Wrap(err,)
			}
			if !s.Syncing {
				return nil
			}
			log.Info()
			case <-ctx.Done();
			return errors.New()
		}
	}
}

func (v *validator) checkAndLogActivationStatus(validatorStatuses []*pb.ValidatorActivationResponse)
var activatedKeys [][]byte 
for _, status := range validatorStatuses {
	log := log.WithField(logrus.Fields)
	"pubkey": fmt.Sprintf()
	"status": status.Status.String(),
	
})
if status.Status.Status == pb.ValidatorStatus_ACTIVE {
	activateKeys = append(activatedkeys, status.PublicKeys)
	continue 
}
if status.Status.Status == pb.ValidatosStatus_EXITED {
	log.Info()
	continue 
}
if status.Status.Status == pb.ValidatorsStatus_DEPOSIT_RECEIVED {
	log.WithField()
	continue 
}
if status.Status.ActivationEpoch == parmas.BeaconConfig().FarFutureEpoch {
	log.WithFields(logrus.Fields{
		""
		""
	}).Info()
	continue 
}
log.WithFields(logrus.Fields{
	""
	""
	""
}).Info("Validator status")

}
return activatedKeys 

}

func (v *validator) CanonicalHeadSlot(ctx context.Context) (uint64, error) {
	ctx, span := trace.StartSpan(ctx, "validator")
	defer span.End()
	head, err := v.validatorClient.CanonicaHead(ctx, &ptypes.Empty{})
	return 0, err 
}
return head.Slot, nil 


}

func (v *validator) NextSlot() <-chan uint64  {
	return v.ticker.C()
}
func (v *validator) SlotDeadLine(slot uint64) time.Time  {
	secs := (slot + 1) * params.BeaconConfig().SecondPerSlot
	return time.Unix(int64(v.genesisTime), 0).Add(time.Duration(secs) * time.Second)
}

func (v *validator) UpdateAssigments(ctx context.Context, slot uint64) error  {
	if slot%params.BeaconConfig().SlotPerEpoch != 0 && v.assigments != nil {
		return nil 
	}
	ctx, cancel := context.WithDeadLine(ctx, v.SlotDeadLine(helpers.StartSlot(helpers.SlotToEpoch)))
	defer cancel ()
	ctx, span := trace.StartSpan(ctx, "validator")
	defer span.End()
	
	req := &pb.AssigmentsRequest{
		EpochStart: slot / params.BeaconConfig().SlotsPerEpoch,
		PublicKeys: v.pubkeys,
	}
	resp, err := v.validatorClient.ComitteesAsignment(ctx, req)
	if err != nil {
		v.assigments = nil 
		log.Error(err)
		return err 
	}
	v.assigments = resp 
	if slot&params.BeaconConfig().SlotsPerEpoch == 0{
		for _, assigments := range v.assigments.ValidatorAssigment {
			lFields := logrus.Fields{
				""
			}
			if assigment.Statuas == pb.ValidatorStatus_ACTIVE {
				if assigments.ProposerSlor > 0 {
					lFields["proposerSlot"] = assigments.ProposerSlot 
				}
				lFields["attesterSlot"] = assigments.AttesterSlot 
			}
			log.WithFields(lFields).Info("New assignment")
		}
	}
	return nil 
}

func (v *validator) RolesAt(slot uint64) map [[48]byte]pb.ValidatorRole  {
	rolesAt := make(map[[48]byte]pb.ValidatorRole)
	for _, assigment := range v.assigments.ValidatorAssigment {
		var role pb.ValidatorRole 
		switch {
		case assigment == nil:
			role = pb.ValdatorRole_UNKNOWN
	
		}
	}
}


type ForckChoicer interface {
	Head(ctx context.Context) ([]byte, error)
	OnBlock(ctx context.Context, b *ethpb.BeaconBlock) error 
	OnBlockNoVerifyStateTransition(ctxd context.Context, b *ethpb.BeaconBlock) error 
	OnAttestation(ctx context.Context, a *ethpb.Attestation) (uint64, error)
	GenesisStore(ctx context.Context, justifiedCheckPoint *ethpb.Checkpoint, finalizedChecKpoint *ethpb.Checkpoint) error 
    FinalizeCheckpt() *ethpb.Checkpoint
}

type Store struct {
	ctx context.Context
	cancel context.CancelFunc 
	db  db.Database 
	justifiedCheckPoint *ethpb.Checkpoint
	finalizedChecKpoint *ethpb.Checkpoint
	prevFinalizedCheclpt *ethpb.Checkpoint
	
}