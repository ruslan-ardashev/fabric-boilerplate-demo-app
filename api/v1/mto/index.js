'use strict';

var express = require('express');
var controller = require('./mto.controller');

var router = express.Router();

router.get('/', controller.mtos);

// maybe expand to include more information about MTOs here

// router.get('/:thingId', controller.getThing);
// router.post('/', controller.addThing);

module.exports = router;
