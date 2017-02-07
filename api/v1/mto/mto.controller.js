'use strict';

const MTO = require('./mto.model');
const BlockchainService = require('../../../services/blockchainSrvc.js');

/*
    Retrieve list of all mtos

    METHOD: GET
    URL : /api/v1/mto
    Response:
        [{'mtoName'}, {'mtoName'}, {'mtoName'}]
*/
exports.mtos = function(req, res) {
    console.log("-- Query list of MTOs --");

    const functionName = "mtos";
    const args = [];

    BlockchainService.query(functionName,args,req.userId).then(function(mtos){
        if (!mtos) {
            res.json([]);
        } else {
            console.log("Retrieved mtos from the blockchain: # " + mtos.length);
            console.log("api-v1-mtos: " + mtos);
            res.json(mtos)
        }
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
};

