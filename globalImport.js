const response = require("cfn-response");
const AWS = require("aws-sdk");
exports.handler = (event, context) => {
  const { SourceRegion } = event.ResourceProperties;
  const id = `custom:cloudformation:${SourceRegion}:exports`;
  const cloudformation = new AWS.CloudFormation({ region: SourceRegion });
  try {
    cloudformation.listExports({}, (err, data) => {
      if (err) {
        throw err;
      } else {
        const obj = {};
        data.Exports.forEach(({ Name, Value }) => (obj[Name] = Value));
        response.send(event, context, response.SUCCESS, obj, id);
      }
    });
  } catch (err) {
    console.error(err.message);
    response.send(event, context, response.FAILED, {}, id);
  }
};
