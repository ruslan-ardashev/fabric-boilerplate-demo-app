'use strict';

const Create = require('./create.model');
const BlockchainService = require('../../../services/blockchainSrvc.js');

/*
    Start create

    METHOD: POST
    URL : /api/v1/transfer
    Arguments:
        // taken directly from chaincode
        //  arg0 - string mto.Name
        //  arg1 - string account.firstName
        //  arg2 - string account.lastName
        //  arg3 - string account.number
        //  arg4 - string account.balance
    Response:
        []
*/
exports.create = function(req, res) {
    console.log("-- Calling Create from api-v1-transfer --");

    var body = req.body;

    console.log("~~req~~");
    console.log(body);
    console.log("~~ENDreq~~");

    const functionName = "createAccount";
    const args = [body.mto, body.firstName, body.lastName, body.accountNumber, body.balance];

    console.log("transfer req: " + req);
    console.log("passed args: " + args);

    BlockchainService.invoke(functionName,args,req.userId).then(function(result){
        if (!result) {
            res.json([]);
        } else {
            // console.log("Retrieved things from the blockchain: # " + result);
            res.json(result)
        }
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
};

/*
    After transfer

    METHOD: POST
    URL : /api/v1/transfer
    Arguments:
        // taken directly from chaincode
        // arg0 - (Source MTO).Name
        // arg1 - (Source account).Number
        // arg2 - balance to transfer from source account
        // arg3 - (Destination MTO).Name
        // arg4 - (Destination account).Number
    Response:
        []
*/
exports.afterCreate = function(req, res) {
    console.log("-- Calling afterCreate from api-v1-transfer-after --");

    var body = req.body;

    console.log("~~req~~");
    console.log(body);
    console.log("~~ENDreq~~");

    const functionName = "afterTransfer";
    const args = [ ];

    console.log("afterTransfer req: " + req);
    console.log("passed args: " + args);

    BlockchainService.query(functionName,args,req.userId).then(function(result){
        if (!result) {
            res.json([]);
        } else {
            // console.log("Retrieved things from the blockchain: # " + result);
            res.json(result)
        }
    }).catch(function(err){
        console.log("Error", err);
        res.sendStatus(500);   
    }); 
};

/*
    Retrieve thing object

    METHOD: GET
    URL: /api/v1/thing/:thingId
    Response:
        { thing }
*/
// exports.getThing = function(req, res) {
//     console.log("-- Query thing --")

//     const functionName = "get_thing";
//     const args = [req.params.thingId];
    
//     BlockchainService.query(functionName,args,req.userId).then(function(thing){
//         if (!thing) {
//             res.json([]);
//         } else {
//             console.log("Retrieved thing from the blockchain");
//             res.json(thing)
//         }
//     }).catch(function(err){
//         console.log("Error", err);
//         res.sendStatus(500);   
//     }); 
// };


//     Add thing object

//     METHOD: POST
//     URL: /api/v1/thing/
//     Response:
//         {  }

// exports.addThing = function(req, res) {
//     console.log("-- Adding thing --")
      
//     const functionName = "add_thing";
//     const args = [req.body.thingId, JSON.stringify(req.body.thing)];
    
//     BlockchainService.invoke(functionName,args,req.userId).then(function(thing){
//         res.sendStatus(200);
//     }).catch(function(err){
//         console.log("Error", err);
//         res.sendStatus(500);   
//     }); 
// };

