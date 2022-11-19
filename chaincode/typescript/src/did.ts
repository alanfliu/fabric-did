/*
 * SPDX-License-Identifier: Apache-2.0
 */

export interface authenticaion {
    kid:string
    method: string
    controller: string
    publicpem: string
}

export interface proof {
    creator: string
    method: string
    signature: any
}

export interface DID {
    id: string
    authenticaions: authenticaion[]
    proof: proof
}
