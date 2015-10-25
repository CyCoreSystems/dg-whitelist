Lists = new Mongo.Collection('lists');

if( Meteor.isClient ) {
   Meteor.subscribe('lists');
}

if( Meteor.isServer ) {
   Meteor.publish('lists',function() {
      return Lists.find({});
   });

   Meteor.startup( function() {
      Lists.find({list: 'whitelist'}).observeChanges({
         'added': function() {
            onUpdate('whitelist');
         },
         'removed': function() {
            onUpdate('whitelist');
         }
      });
      Lists.find({list: 'greylist'}).observeChanges({
         'added': function() {
            onUpdate('greylist');
         },
         'removed': function() {
            onUpdate('greylist');
         }
      });
      Lists.find({list: 'blacklist'}).observeChanges({
         'added': function() {
            onUpdate('blacklist');
         },
         'removed': function() {
            onUpdate('blacklist');
         }
      });
   });
}

var onUpdate = function(list) {
   dg.export(list, Lists.find({ list: list }));
};
