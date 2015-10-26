Lists = new Mongo.Collection('lists');

if( Meteor.isClient ) {
   Meteor.subscribe('lists');
}

if( Meteor.isServer ) {
   Meteor.publish('lists',function() {
      return Lists.find({});
   });
}
