
import neo4j from 'k6/x/neo4j'

export default function () {
    const neo4jConf = {
        uri: "bolt://localhost:7687",
        user: "neo4j",
        password: "neo4jpass",
        realm: ""
    }
    // instantiate the driver
    const driver = neo4j.newDriver(neo4jConf)

    // execute write query
    const params = {
        name: "Errico",
        country: "Italia",
        year: 1853,
    }
    // function write(query string, params object)
    const writeResult = driver.write("CREATE (p:Person {name: $name, country: $country, year: $year}) RETURN p;", params)
    console.log(writeResult) // [{"values":[{"id":0,"element_id":"4:17dbeda4-05f9-46b3-8fe2-9c480144afda:0","labels":["Person"],"props":{"year":1853,"name":"Errico","country":"Italia"}}],"keys":["p"]}]

    // execute read query
    // function read(query string, params object)
    const readResult = driver.read("MATCH (p:Person) WHERE p.name = 'Errico' RETURN p;")
    console.log(readResult) // [{"values":[{"id":0,"element_id":"4:17dbeda4-05f9-46b3-8fe2-9c480144afda:0","labels":["Person"],"props":{"country":"Italia","year":1853,"name":"Errico"}}],"keys":["p"]}]


    // execute query
    // read and write functions call this function and define the access mode
    // you can use this function directly and set the access mode yourself
    // function executeQuery(accessMode int, query string, params object)
    // valid access mode values: [0, 1]
    //   0 - Write: tells the driver to use a connection to 'Leader'
    //   1 - Read: tells the driver to use a connection to one of the 'Follower' or 'Read Replica'.
    const writeExecResult = driver.executeQuery(0, "CREATE (p:Person {name: 'Murray', country: 'USA', year: 1921}) RETURN p;", {
        name: "Murray",
        country: 'America',
        year: 1921})
    console.log(writeExecResult) // [{"values":[{"id":1,"element_id":"4:17dbeda4-05f9-46b3-8fe2-9c480144afda:1","labels":["Person"],"props":{"country":"USA","year":1921,"name":"Murray"}}],"keys":["p"]}]

    // Can use without parameters
    const readExecResult = driver.executeQuery(1, "MATCH (p:Person) WHERE p.name = 'Murray' RETURN p;")
    console.log(readExecResult) // [{"values":[{"id":1,"element_id":"4:17dbeda4-05f9-46b3-8fe2-9c480144afda:1","labels":["Person"],"props":{"country":"USA","year":1921,"name":"Murray"}}],"keys":["p"]}]
}