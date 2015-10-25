ListTypes = new Mongo.Collection('listTypes');

if( Meteor.isClient ) {
   Meteor.subscribe('listTypes');
}

if( Meteor.isServer ) {
   ListTypes.remove({});
   ListTypes.insert({
      name: 'whitelist',
      label: 'Whitelist',
      iconName: 'done',
      order: 0,
   });
   ListTypes.insert({
      name: 'greylist',
      label: 'Greylist',
      iconName: 'call_split',
      order: 0,
   });
   ListTypes.insert({
      name: 'blacklist',
      label: 'Blacklist',
      iconName: 'not_interested',
      order: 0,
   });
   Meteor.publish('listTypes',function() {
      return ListTypes.find({});
   });
}
