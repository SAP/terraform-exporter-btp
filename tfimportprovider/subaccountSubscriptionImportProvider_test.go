package tfimportprovider

// 	jsonString := "{\"id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\",\"subaccount_id\":\"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5\",\"values\":[{\"state\":\"SUBSCRIBED\",\"app_name\":\"testapp\",\"plan_name\":\"testplan\"}]}"
// 	expectedValue := "import {\n\t\t\t\tto = btp_subaccount_subscription.testapp\n\t\t\t\tid = \"5163621f-6a1e-4fbf-af3a-0f530a0dc4d5,testapp,testplan\"\n\t\t\t  }\n"
