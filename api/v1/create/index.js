'use strict';

var express = require('express');
var controller = require('./create.controller');

var router = express.Router();

router.post('/', controller.create);
router.get('/after', controller.afterCreate);
// router.get('/:thingId', controller.getThing);
// router.post('/', controller.addThing);

module.exports = router;
