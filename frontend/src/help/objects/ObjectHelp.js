
export default {
  title: "What are objects?",
  paragraphs: [
    `Objects describe the type of data stored in your database. Each object represents a different data type. Objects use attributes to specify 
    the information stored for that data type. For example, suppose you are creating an app that allows you to view the dogs at different animal 
    shelters. You would probably want to have a "Dog" object. The Dog object could have attributes such as "name", "breed" and "age." CRUD Creator
    then uses your object to create a database that can store those values.`,

    `Note that the Dog object does not represent a specific, real-life dog. Rather, it represents the data stored for dogs in general. We create 
    an attribute called "breed", not an attribute called "beagle". "Beagle" would later be stored in your database as the value for "breed" when 
    you add a specific, real-life dog to your database. When defining your own objects, you want to think in general terms like this. Create 
    specific instances of your objects after your database is already deployed.`
  ]
}
