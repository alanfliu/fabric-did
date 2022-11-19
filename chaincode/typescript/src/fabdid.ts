/*
 * SPDX-License-Identifier: Apache-2.0
 */

import { Context, Contract } from 'fabric-contract-api';
import { DID } from './did';

export class FabDID extends Contract {

    public async initLedger(ctx: Context) {
        console.info('============= START : Initialize Ledger ===========');
        console.info('============= END : Initialize Ledger ===========');
    }

    public async queryDocument(ctx: Context, id: string): Promise<string> {
        const didAsBytes = await ctx.stub.getState(id); // get the car from chaincode state
        if (!didAsBytes || didAsBytes.length === 0) {
            throw new Error(`${id} does not exist`);
        }
        console.log(didAsBytes.toString());
        return didAsBytes.toString();
    }

    public async createCar(ctx: Context, didjson: string) {
        const did: DID = JSON.parse(didjson)
        if(!did || !did.id) {
            throw new Error("parse did json err.")
        }

        const exit = await ctx.stub.getState(did.id)
        if(exit || exit.length != 0) {
            throw new Error(`${did.id} already exit`)
        }

        await ctx.stub.putState(did.id, Buffer.from(JSON.stringify(did)))
    }
}
