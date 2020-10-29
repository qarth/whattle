package main

import "github.com/apex/log"

func DryRun(opt RunCtx) error {

	log.Infof("Begin parsing parameters from %v", opt.ParamFile)
	notifyStatus(opt.Notify, "Parsing parameters file")

	var params Parameters

	if e := readJSONFile(opt.ParamFile, &params); e != nil {
		return e
	}

	log.Infof("Begin reading input from %v", opt.InputFile)
	notifyStatus(opt.Notify, "Reading input data")

	if e := params.Input.initializeFromGzip(opt.InputFile); e != nil {
		return e
	}

	nReal := len(ctx.Input.Ebv)
	nData := len(ctx.Input.Ebv[0])
	log.Infof("Number of realizations: %v", nReal)
	log.Infof("Number of rows: %v", nData)
	log.Info("Begin creating naive mask")
	//notifyStatus(ch, "Creating naive mask")
	mask := ctx.generateMask()
	log.Info("Begin creating precedence")
	//notifyStatus(ch, "Creating precedence")
	if ctx.Precedence.init(ctx, mask) != nil {
		return nil, -1
	}
	log.Info("Updating mask")
	//notifyStatus(ch, "Updating mask")
	for i := 0; i < nData; i++ {
		if mask[i] {
			if key := ctx.Precedence.keys[i]; key != MISSING {
				for _, off := range ctx.Precedence.defs[key] {
					mask[i+off] = true
				}
			}
		}
	}

	log.Info("Begin compressing")
	//notifyStatus(ch, "Compressing")
	var condensedEBV Data
	var condensedPre Precedence

	if c := compressEverything(mask, &ctx.Input, &ctx.Precedence, &condensedEBV, &condensedPre); err != nil {
		log.Info("ERROR: Compressing everything failed")
		return nil, 1
	} else {
		log.Info("SUCCESS: Compressed everything")
	}
}
