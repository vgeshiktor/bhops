import asyncio
import configparser

from msgraph.generated.models.o_data_errors.o_data_error import ODataError

from graph import Graph


async def main():
    print('Python Graph App-Only Tutorial\n')

    # Load settings
    config = configparser.ConfigParser()
    config.read(['./config.cfg', 'config.dev.cfg'])
    azure_settings = config['azure']

    graph: Graph = Graph(azure_settings)

    choice = -1

    while choice != 0:
        print('Please choose one of the following options:')
        print('0. Exit')
        print('1. Display access token')
        print('2. List users')
        print('3. Make a Graph call')

        try:
            choice = int(input())
        except ValueError:
            choice = -1

        try:
            if choice == 0:
                print('Goodbye...')
            elif choice == 1:
                await display_access_token(graph)
            elif choice == 2:
                await list_users(graph)
            elif choice == 3:
                await make_graph_call(graph)
            else:
                print('Invalid choice!\n')
        except ODataError as odata_error:
            print('Error:')
            if odata_error.error:
                print(odata_error.error.code, odata_error.error.message)


async def display_access_token(graph: Graph):
    # TODO
    return


async def list_users(graph: Graph):
    # TODO
    return


async def make_graph_call(graph: Graph):
    # TODO
    return


# Run main
asyncio.run(main())