'use strict';

var express = require('express');
var controller = require('./transfer.controller');

var router = express.Router();

router.post('/', controller.transfer);
router.get('/after', controller.afterTransfer);
// router.get('/:thingId', controller.getThing);
// router.post('/', controller.addThing);

module.exports = router;
