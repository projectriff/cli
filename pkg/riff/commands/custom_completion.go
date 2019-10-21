package commands

const bash_completion_func = `
__riff_list_resource()
{
    local riff_output out
    if riff_output=$(riff $1 2>/dev/null); then
        out=($(echo "${riff_output}" | awk 'NR>1 {print $1}'))
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__riff_custom_func() {
    case ${last_command} in
        riff_application_delete | riff_application_status | riff_application_tail)
            __riff_list_resource 'application list'
            return
            ;;
        riff_container_delete | riff_container_status)
            __riff_list_resource 'container list'
            return
            ;;
        riff_core_deployer_delete | riff_core_deployer_status | riff_core_deployer_tail)
            __riff_list_resource 'core deployer list'
            return
            ;;
        riff_crediential_delete)
            __riff_list_resource 'credential list'
            return
            ;;
        riff_function_delete | riff_function_status | riff_function_tail)
            __riff_list_resource 'function list'
            return
            ;;
        riff_knative_deployer_delete | riff_knative_deployer_status | riff_knative_deployer_tail)
            __riff_list_resource 'knative deployer list'
            return
            ;;
        riff_knative_adapter_delete | riff_knative_adapter_status)
            __riff_list_resource 'knative adapter list'
            return
            ;;
        riff_streaming_kafka-provider_delete | riff_streaming_kafka-provider_status)
            __riff_list_resource 'streaming kafka-provider list'
            return
            ;;
        riff_streaming_processor_delete | riff_streaming_processor_status | riff_streaming_processor_tail)
            __riff_list_resource 'streaming processor list'
            return
            ;;
        riff_streaming_stream_delete | riff_streaming_stream_status)
            __riff_list_resource 'streaming stream list'
            return
            ;;
        *)
            ;;
    esac
}
`
